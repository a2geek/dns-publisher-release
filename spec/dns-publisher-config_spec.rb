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

      expect(config['Trigger']['Type']).to eq("file-watcher")
      expect(config['Trigger']['FileWatcher']).to eq("/var/vcap/instance/dns/records.json")
    end

    it 'allows file-watcher trigger overrides' do
      config = JSON.parse(template.render({
        "trigger" => {
          "file-watcher" => "/tmp/file.json"
        }
      }))

      expect(config['Trigger']['Type']).to eq("file-watcher")
      expect(config['Trigger']['FileWatcher']).to eq("/tmp/file.json")
    end

    it 'allows timer trigger selection' do
      config = JSON.parse(template.render({
        "trigger" => {
          "type" => "timer",
          "refresh" => "30s"
        }
      }))

      expect(config['Trigger']['Type']).to eq("timer")
      expect(config['Trigger']['Refresh']).to eq("30s")
    end

    it 'assigns mappings defaults' do
      config = JSON.parse(template.render({
          "mappings" => [
            { 
              "instance-group" => "concourse",
              "deployment" => "concourse",
              "fqdn" => "concourse.lan"
            }
          ]
        }))

        expect(config["Mappings"][0]["Network"]).to eq("default")
        expect(config["Mappings"][0]["TLD"]).to eq("bosh")
    end

    it 'allows default overrides in mappings' do
      config = JSON.parse(template.render({
          "mappings" => [
            { 
              "instance-group" => "concourse",
              "network" => "my-network",
              "deployment" => "concourse",
              "tld" => "my-tld",
              "fqdn" => "concourse.lan"
            }
          ]
        }))

        expect(config["Mappings"][0]["Network"]).to eq("my-network")
        expect(config["Mappings"][0]["TLD"]).to eq("my-tld")
    end
  end
end
