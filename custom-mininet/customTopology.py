import time
from random import randint

from mininet.net import Mininet
from mininet.topo import Topo
from mininet.node import Controller, OVSSwitch
from mininet.cli import CLI
from mininet.log import setLogLevel, info, output
from mininet.link import TCLink

class MobilitySwitch( OVSSwitch ):
    "Switch that can reattach and rename interfaces"

    def delIntf( self, intf ):
        "Remove (and detach) an interface"
        port = self.ports[ intf ]
        del self.ports[ intf ]
        del self.intfs[ port ]
        del self.nameToIntf[ intf.name ]

    # pylint: disable=arguments-differ
    def addIntf( self, intf, rename=False, **kwargs ):
        "Add (and reparent) an interface"
        OVSSwitch.addIntf( self, intf, **kwargs )
        intf.node = self
        if rename:
            self.renameIntf( intf )

    def attach( self, intf ):
        "Attach an interface and set its port"
        port = self.ports[ intf ]
        if port:
            if self.isOldOVS():
                self.cmd( 'ovs-vsctl add-port', self, intf )
            else:
                self.cmd( 'ovs-vsctl add-port', self, intf,
                          '-- set Interface', intf,
                          'ofport_request=%s' % port )
            self.validatePort( intf )

    def validatePort( self, intf ):
        "Validate intf's OF port number"
        ofport = int( self.cmd( 'ovs-vsctl get Interface', intf,
                                'ofport' ) )
        if ofport != self.ports[ intf ]:
            warn( 'WARNING: ofport for', intf, 'is actually', ofport,
                  '\n' )

    def renameIntf( self, intf, newname='' ):
        "Rename an interface (to its canonical name)"
        intf.ifconfig( 'down' )
        if not newname:
            newname = '%s-eth%d' % ( self.name, self.ports[ intf ] )
        intf.cmd( 'ip link set', intf, 'name', newname )
        del self.nameToIntf[ intf.name ]
        intf.name = newname
        self.nameToIntf[ intf.name ] = intf
        intf.ifconfig( 'up' )

    def moveIntf( self, intf, switch, port=None, rename=True ):
        "Move one of our interfaces to another switch"
        self.detach( intf )
        self.delIntf( intf )
        switch.addIntf( intf, port=port, rename=rename )
        switch.attach( intf )

def printConnections( switches ):
    "Compactly print connected nodes to each switch"
    for sw in switches:
        output( '%s: ' % sw )
        for intf in sw.intfList():
            link = intf.link
            if link:
                intf1, intf2 = link.intf1, link.intf2
                remote = intf1 if intf1.node != sw else intf2
                output( '%s(%s) ' % ( remote.node, sw.ports[ intf ] ) )
        output( '\n' )

def moveHost( host, oldSwitch, newSwitch, newPort=None ):
    "Move a host from old switch to new switch"
    hintf, sintf = host.connectionsTo( oldSwitch )[ 0 ]
    oldSwitch.moveIntf( sintf, newSwitch, port=newPort )
    return hintf, sintf


#Topologia Custom
class NetworkTopo(Topo):
	def build(self, **_opts):
		
		#Switch: [s1, s2]
		s1 = self.addSwitch('s1')
		s2 = self.addSwitch('s2')
		s3 = self.addSwitch('s3')
		#s4 = self.addSwitch('s4')
		#s5 = self.addSwitch('s5')
		#s6 = self.addSwitch('s6')

		#Host: [h1, h2, h3]
		h1 = self.addHost('h1', ip = '10.0.0.1/28')
		h2 = self.addHost('h2', ip = '10.0.0.2/28')
		h3 = self.addHost('h3', ip = '10.0.0.3/28')

		#[s1: h1, h2]
                self.addLink(h1, s3)
                self.addLink(h2, s1)
		self.addLink(h3, s3)

		#Collegamenti:
		#[s1 --- s2]
     		self.addLink(s2, s1) #, bw=0.3, delay='100ms', loss=3, max_queue_size=1000, use_htb=True)
		self.addLink(s2, s3) #, bw=0.3, delay='100ms', loss=3, max_queue_size=1000, use_htb=True)
		#self.addLink(s3, s4, bw=2, delay='50ms', loss=3, max_queue_size=1000, use_htb=True)
		#self.addLink(s4, s5, bw=2, delay='50ms', loss=3, max_queue_size=1000, use_htb=True)
		#self.addLink(s4, s6, bw=2, delay='50ms', loss=3, max_queue_size=1000, use_htb=True)


def test_ping(host1, host2):
    result = host1.cmd('ping -c 3 %s' % host2.IP())
    print(result)




def run():
	
	net = Mininet(topo=NetworkTopo(), switch=MobilitySwitch, link=TCLink)
	
	#Start network
	net.start()

	printConnections( net.switches )

	#Nodes
	h1, h2, s1, s2 = net.get('h1', 'h2', 's1', 's2')
	h1.cmd('cd /root/quic-go/example')
	#h1.cmd('go run serverTCP.go &')
	#h1.cmd('go run main.go &')
	h1.cmd('go run cont.go &')

	#Test ping before mobility
	print("*** Before mobility")
	#test_ping(h2, h1)
	#net.pingAll()
	time.sleep(2)

	h2.cmd('cd /root/quic-go/example/client')
	#h2.cmd('go run clientTCP.go https://10.0.0.1:8080/page &')
	#h2.cmd('go run main.go -insecure https://10.0.0.1:6121/page &')
	h2.cmd('go run cont.go &')	
	
	time.sleep(10)

	#Mobility
	h1, old = net.get( 'h1', 's3' )
        new = net[ 's1' ]
	port = randint( 10, 20 )
        info( '* Moving', h1, 'from', old, 'to', new, 'port', port, '\n' )
        hintf, sintf = moveHost( h1, old, new, newPort=port )
        info( '*', hintf, 'is now connected to', sintf, '\n' )
        info( '* Clearing out old flows\n' )
        for sw in net.switches:
            sw.dpctl( 'del-flows' )
        info( '* New network:\n' )
        printConnections( net.switches )
	

	#Test ping after mobility
	print("*** After mobility")
        #test_ping(h2, h1)
	#net.pingAll()

	CLI(net)
	net.stop()

if __name__ == '__main__':
	setLogLevel('info')
	run()
