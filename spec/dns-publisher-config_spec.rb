require 'rspec'
require 'json'
require 'bosh/template/test'

describe 'dns-publisher job' do
  let(:release) { Bosh::Template::Test::ReleaseDir.new(File.join(File.dirname(__FILE__), '..')) }
  let(:job) { release.job('dns-publisher') }

  describe 'config.json' do
    let(:template) { job.template('config/config.json') }

    it 'assigns trigger defaults' do
      config = JSON.parse(template.render({
        # all defaults
      }))

      expect(config['BoshDns']['Trigger']['Type']).to eq("file-watcher")
      expect(config['BoshDns']['Trigger']['FileWatcher']).to eq("/var/vcap/instance/dns/records.json")
    end

    it 'allows file-watcher trigger overrides' do
      config = JSON.parse(template.render({
        "bosh-dns" => {
            "trigger" => {
              "file-watcher" => "/tmp/file.json"
            }
          }
      }))

      expect(config['BoshDns']['Trigger']['Type']).to eq("file-watcher")
      expect(config['BoshDns']['Trigger']['FileWatcher']).to eq("/tmp/file.json")
    end

    it 'allows timer trigger selection' do
      config = JSON.parse(template.render({
        "bosh-dns" => {
          "trigger" => {
            "type" => "timer",
            "refresh" => "30s"
          }
        }
      }))

      expect(config['BoshDns']['Trigger']['Type']).to eq("timer")
      expect(config['BoshDns']['Trigger']['Refresh']).to eq("30s")
    end

    it 'assigns mappings defaults' do
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
    end

    it 'allows default overrides in mappings' do
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
    end
  end
end
