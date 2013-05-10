metadata    :name        => "fileaudit",
            :description => "Return MD5sum for a file",
            :author      => "Chris Mague",
            :license     => "Apache License 2.0",
            :version     => "0.1",
            :url         => "http://blog.mague.com",
            :timeout     => 20

action "audit", :description => "Return the md5sum for a file" do
    display :always

    output :output,
           :description => "MD5sum output",
           :display_as => "md5sum"
end

action "role_call", :description => "check existance of a file" do
    display :always

    output :output,
           :description => "status",
           :display_as => "status"
end

action "tfe", :description => "Test for Echo" do
    display :always

    output :status,
           :description => "Status Message",
           :display_as => "Status"
end

action "get_link", :description => "check existance of a file" do
    display :always

    output :output,
           :description => "status",
           :display_as => "status"
end
