require 'rspec'
require 'json'
require 'bosh/template/test'

describe 'dns-publisher job' do
  let(:release) { Bosh::Template::Test::ReleaseDir.new(File.join(File.dirname(__FILE__), '..')) }
  let(:job) { release.job('dns-publisher') }

  describe 'config.json' do
    let(:template) { job.template('config/config.json') }

    it 'assigns bosh dns trigger defaults' do
      config = JSON.parse(template.render({
        # all defaults
        "bosh-dns" => {
          "mappings" => []
        }
      }))

      expect(config['BoshDns']['Trigger']['Type']).to eq("file-watcher")
      expect(config['BoshDns']['Trigger']['FileWatcher']).to eq("/var/vcap/instance/dns/records.json")
      expect(config['CloudFoundry']).to be(nil)
    end

    it 'allows bosh dns file-watcher trigger overrides' do
      config = JSON.parse(template.render({
        "bosh-dns" => {
            "trigger" => {
              "file-watcher" => "/tmp/file.json"
            },
            "mappings" => []
          }
      }))

      expect(config['BoshDns']['Trigger']['Type']).to eq("file-watcher")
      expect(config['BoshDns']['Trigger']['FileWatcher']).to eq("/tmp/file.json")
      expect(config['CloudFoundry']).to be(nil)
    end

    it 'allows bosh dns timer trigger selection' do
      config = JSON.parse(template.render({
        "bosh-dns" => {
          "trigger" => {
            "type" => "timer",
            "refresh" => "30s"
          },
          "mappings" => []
        }
      }))

      expect(config['BoshDns']['Trigger']['Type']).to eq("timer")
      expect(config['BoshDns']['Trigger']['Refresh']).to eq("30s")
      expect(config['CloudFoundry']).to be(nil)
    end

    it 'assigns bosh dns mappings defaults' do
      config = JSON.parse(template.render({
        "bosh-dns" => {
          "mappings" => [
            { 
              "instance-group" => "concourse",
              "deployment" => "concourse",
              "fqdns" => ["concourse.lan"]
            }
          ]
        }
      }))

      expect(config['BoshDns']["Mappings"][0]["Network"]).to eq("default")
      expect(config['BoshDns']["Mappings"][0]["TLD"]).to eq("bosh")
      expect(config['CloudFoundry']).to be(nil)
    end

    it 'allows bosh dns default overrides in mappings' do
      config = JSON.parse(template.render({
        "bosh-dns" => {
          "mappings" => [
            { 
              "instance-group" => "concourse",
              "network" => "my-network",
              "deployment" => "concourse",
              "tld" => "my-tld",
              "fqdns" => ["concourse.lan"]
            }
          ]
        }
      }))

      expect(config['BoshDns']["Mappings"][0]["Network"]).to eq("my-network")
      expect(config['BoshDns']["Mappings"][0]["TLD"]).to eq("my-tld")
      expect(config['CloudFoundry']).to eq(nil)
    end

    it 'allows cloud foundry defaults' do
      config = JSON.parse(template.render({
        "cloud-foundry" => {
          "url" => "https://api.sys.cf.lan",
          "client-id" => "some_client_id",
          "client-secret" => "sekret",
          "alias" => "my.alias.com",
          "mappings" => ["*.cf.lan", "*.myapp.lan"]
        }
      }))

      expect(config['BoshDns']).to eq(nil)
      expect(config['CloudFoundry']["Trigger"]["Type"]).to eq("timer")
      expect(config['CloudFoundry']["Trigger"]["Refresh"]).to eq("5m")
      expect(config['CloudFoundry']["URL"]).to eq("https://api.sys.cf.lan")
      expect(config['CloudFoundry']["SkipSslValidation"]).to eq(false)
      expect(config['CloudFoundry']["ClientId"]).to eq("some_client_id")
      expect(config['CloudFoundry']["ClientSecret"]).to eq("sekret")
      expect(config['CloudFoundry']["Alias"]).to eq("my.alias.com")
      expect(config['CloudFoundry']["Mappings"]).to match_array(["*.cf.lan", "*.myapp.lan"])
    end
  end
end
