#!/usr/bin/env ruby
################################################################################
# Usage
#
#   print out the resources that are available
# 
#  ./consul_remote_state_outputs.rb --cluster consul-stage.example.com --path terraform/stage/core_network
#  {"subnets-app"=>
#    {"sensitive"=>false,
#     "type"=>"list",
#     "value"=>["subnet-8c8941e5", "subnet-6bbfa413"]},
#   "subnets-management"=>
#    {"sensitive"=>false,
#     "type"=>"list",
#     "value"=>["subnet-8a8941e3", "subnet-6cbfa414"]},
#   "subnets-public"=>
#    {"sensitive"=>false,
#     "type"=>"list",
#     "value"=>["subnet-8d8941e4", "subnet-6abfa412"]},
#   "vpc-id"=>{"sensitive"=>false, "type"=>"string", "value"=>"vpc-97c138fe"}}


require 'getoptlong'
require 'uri'
require 'json'
require 'net/http'
require 'base64'
require 'pp'


args = {  
          :cluster => 'localhost',
          :port    => 8500,
          :path    => '/',
}

opts = GetoptLong.new(
  [ '--cluster'        , '-c' , GetoptLong::OPTIONAL_ARGUMENT ],
  [ '--port'           , '-p' , GetoptLong::OPTIONAL_ARGUMENT ],
  [ '--path'           , '-P' , GetoptLong::OPTIONAL_ARGUMENT ],
)

opts.each do |opt, arg|
  case opt
    when '--cluster'
      args[:cluster] = arg
    when '--path'
      args[:path] = arg
    when '--port'
      args[:port] = arg.to_i
  end
end


begin
  url = URI.parse "http://#{args[:cluster]}:#{args[:port]}/v1/kv/#{args[:path]}"
  #JSON parse the docoded portion of the JSON response 
  JSON.parse(Base64.decode64(JSON.parse(Net::HTTP.get(url)).first['Value']))['modules'].each do |x|
    # In terraform you can only use the root path output data 
    if x['path'] == ['root']
      pp x['outputs']
    end
  end
rescue Exception => e
  STDERR.puts "Error retrieving: #{url}: #{e.message}"
  exit! 1
end
