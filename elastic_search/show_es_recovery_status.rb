#!/bin/env ruby
#################################################################
#   Prints out information on what is happening when ElasticSearch
#   is recovering.  Run in a loop to see what's happening
#   while true; do ./show_es_recovery_status.rb es-master; done
#################################################################

require 'rubygems'
require 'json'
require 'net/http'
require 'uri'
require 'pp'

if ARGV[0] == ""
  puts "run with hostname as first argument"
  exit! 1
end

BASE_URL = "http://#{ARGV[0]}:9200"

def return_info(path)
  es_url = URI.parse("#{BASE_URL}#{path}")
  response = Net::HTTP.get_response(es_url)
  JSON.parse(response.body)
end

nodes = return_info("/_nodes")

puts "#######################################################################"
x = return_info("/_cluster/health")
x.keys.sort.each do |k|
  printf("%30s  %s\n", k, x[k])
end
puts "#######################################################################"
status = return_info("/_status")
status['indices'].each do |index, info|
  info['shards'].each do |shard, vals|
    if vals[0]['state'] == "RECOVERING"
      puts "Index: #{index}  Shard: #{shard} Server: #{nodes['nodes'][vals[0]['routing']['node']]['name']} Size: #{vals[0]['index']['size']}"
    end
  end
end
puts "#######################################################################"