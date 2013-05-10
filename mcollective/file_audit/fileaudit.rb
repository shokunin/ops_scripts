module MCollective
  module Agent

    require 'fileutils'
    require 'digest/md5'

    #############################################################
    # An agent to audit files
    #
    # Configuration Options:
    #    none
    #                        
    #############################################################
    class Fileaudit<RPC::Agent
      metadata    :name        => "fileaudit",
                  :description => "Return md5sum of a file",
                  :author      => "Chris Mague",
                  :license     => "Apache License 2.0",
                  :version     => "0.1",
                  :url         => "http://blog.mague.com",
                  :timeout     => 30

      #############################################################
      action "audit" do
        if File.exists? request[:filename]
          reply[:output] = audit_file(request[:filename])
        else
          reply[:output] = '0'
        end
      end

      action "role_call" do
        if File.exists? request[:filename]
          reply[:output] = "PRESENT"
        else
          reply[:output] = "ABSENT"
        end
      end

      action "get_link" do
        if File.exists? request[:filename] and File.symlink? request[:filename]
          reply[:output] = File.readlink request[:filename]
        else
          reply[:output] = "ERROR"
        end
      end

      def audit_file(filename)
        begin
          Digest::MD5.hexdigest(File.read(filename))
        rescue Exception => e
          reply.fail!("Fail: #{e.message}")
        end
      end
    end
  end
end
