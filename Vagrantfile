# -*- mode: ruby -*-
# vi: set ft=ruby :

# All Vagrant configuration is done below. The "2" in Vagrant.configure
# configures the configuration version (we support older styles for
# backwards compatibility). Please don't change it unless you know what
# you're doing.
Vagrant.configure("2") do |config|

  config.vm.box = "base"

  config.vm.provision "shell",
    inline: "yum -y install golang"

  config.vm.provision "shell",
    inline: "go build -o /usr/bin/congruit /vagrant/congruit.go"

  config.vm.hostname = 'centos7'
  config.vm.box = "geerlingguy/centos7"

  if ENV['WORKPLACES_ENABLED']
    config.vm.provision "shell",
      inline: "yum install git -y && congruit -debug -stockroom-dir=/vagrant/stockroom/ -debug -gitrepo https://github.com/Congruit/example-stockroom.git -workplaces #{ENV['WORKPLACES_ENABLED']}"

    config.vm.define 'Centos7' do |centos7|

    end
  end

$script = <<SCRIPT
pkill  congruit || echo 'congruit is not running'
cd /vagrant
congruit -debug -stockroom-dir=/vagrant/stockroom/ -friend -token foobar -debug -gitrepo https://github.com/Congruit/example-stockroom.git -workplaces tomcat-docker -ssl_cert /vagrant/insecure-domain.crt -ssl_key /vagrant/insecure-domain.key
SCRIPT

  config.vm.define 'Docker01' do |docker01|
    docker01.vm.network "private_network", ip: "192.168.50.4"
    docker01.vm.hostname = 'docker01'
    docker01.vm.provision "shell",
      inline: $script
  end

  config.vm.define 'Docker02' do |docker02|
    docker02.vm.network "private_network", ip: "192.168.50.5"
    docker02.vm.hostname = 'docker02'
    docker02.vm.provision "shell",
      inline: $script
  end

end
