#!/usr/bin/env ruby
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
  begin
    es_url = URI.parse("#{BASE_URL}#{path}")
    http = Net::HTTP.new(es_url.host, es_url.port)
    http.read_timeout = 500
    response = http.request(Net::HTTP::Get.new(es_url.request_uri))
    JSON.parse(response.body)
  rescue Exception => e
    puts "ERROR: #{e.message}"
  end
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
    vals.each do |v|
      if v['routing']['state'] == "RELOCATING"
        begin
          recover_percent=return_info("/#{index}/_recovery?active_only=true")
          per = recover_percent[index]['shards'].select{ |z| z['id'].to_s == shard }.map{ |rshard|  rshard['index']['files']['percent']}[0]
          puts "RELOCATING: #{per} Index: #{v['routing']['index']} Shard: #{v['routing']['shard']} from #{nodes['nodes'][v['routing']['node']]['name']} to #{nodes['nodes'][v['routing']['relocating_node']]['name']} Size: #{v['index']['size_in_bytes']} "
        rescue
          puts "RELOCATING: Index: #{v['routing']['index']} Shard: #{v['routing']['shard']} from #{nodes['nodes'][v['routing']['node']]['name']} to #{nodes['nodes'][v['routing']['relocating_node']]['name']} Size: #{v['index']['size_in_bytes']} "
        end
      end
    end
    if vals[0]['state'] == "RECOVERING"
      begin
        recover_percent=return_info("/#{index}/_recovery?active_only=true")
        per = recover_percent[index]['shards'].select{ |z| z['id'].to_s == shard }.map{ |rshard|  rshard['index']['files']['percent']}[0]
        puts "RECOVERING: #{per} #{vals[0]['state']} Index: #{index}  Shard: #{shard} Server: #{nodes['nodes'][vals[0]['routing']['node']]['name']} Size: #{vals[0]['index']['size_in_bytes']}"
      rescue
        puts "RECOVERING: #{vals[0]['state']} Index: #{index}  Shard: #{shard} Server: #{nodes['nodes'][vals[0]['routing']['node']]['name']} Size: #{vals[0]['index']['size_in_bytes']}"
      end
    end
  end
end
puts "#######################################################################"
