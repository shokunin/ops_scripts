#!/usr/bin/env ruby

require 'rubygems'
require 'net/http'
require 'erb'
require 'uri'
require 'json'
require 'getoptlong'
require 'base64'

args = {
    :host       => 'localhost',
    :port       => 8500,
    :path       => '',
    :datacenter => '',
}

opts = GetoptLong.new(
  [ '--host', '-H', GetoptLong::OPTIONAL_ARGUMENT ],
  [ '--port', '-p', GetoptLong::OPTIONAL_ARGUMENT ],
  [ '--datacenter', '-d', GetoptLong::OPTIONAL_ARGUMENT ],
  [ '--path', '-P', GetoptLong::OPTIONAL_ARGUMENT ]
)

opts.each do |opt, arg|
  case opt
    when '--host'
      args[:host] = arg
    when '--port'
      args[:port] = arg.to_i
    when '--path'
      args[:path] = arg
    when '--datacenter'
      args[:datacenter] = arg
  end
end

def api_fetch(url, params)
  begin
    url = URI.parse(url)
    url.query = URI.encode_www_form( params )
    JSON.parse(Net::HTTP.get(url))
  rescue Exception => e
    STDERR.puts "Error retrieving: #{url}: #{e.message}"
    exit! 1
  end
end

if args[:path] == "" || args[:datacenter] == ""
  puts "need a path and a datacenter to run"
  exit! 1
end

top_template = ERB.new <<-EOF
provider "consul" {
  address    = "<%= args[:host] %>:<%= args[:port] %>"
  datacenter = "<%= args[:datacenter] %>"
  scheme     = "http"
}

resource "consul_keys" "settings" {
  datacenter = "<%= args[:datacenter] %>"
EOF

puts top_template.result(binding)

key_template = ERB.new <<-EOF2
  key {
    name  = "<%= File.basename(k['Key']) %>"
    path  = "<%= k['Key'] %>"
    value = "<%= Base64.decode64(k['Value']) %>"
  }
EOF2

api_fetch("http://#{args[:host]}:#{args[:port]}/v1/kv/#{args[:path]}", {"recurse" => '1'}).each do |k|
  puts key_template.result(binding)
end

puts "}"
