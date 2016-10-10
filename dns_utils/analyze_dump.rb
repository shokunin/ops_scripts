#!/usr/bin/env ruby

require 'pcap'
require 'dnsruby'

begin 
  inp = Pcap::Capture.open_offline(ARGV[0])
  inp.loop(-1) do |pkt|
    data = pkt.udp_data
    #pp Dnsruby::Message.decode(data)
  end
rescue Exception => e
  STDERR.puts e.message
end
