require 'rspec'
require 'json'
require 'bosh/template/test'

describe 'dns-publisher job' do
  let(:release) { Bosh::Template::Test::ReleaseDir.new(File.join(File.dirname(__FILE__), '..')) }
  let(:job) { release.job('dns-publisher') }

  describe 'config.json' do
    let(:template) { job.template('config/config.json') }

    it 'assigns mappings defaults' do
      config = JSON.parse(template.render({
          "mappings" => [
            { 
              "instance-group" => "concourse",
              "deployment" => "concourse",
              "fqdn" => "concourse.lan"
            }
          ],
          "publisher" => {
            "options" => {
              "user" => "root"
            }
          }
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
          ],
          "publisher" => {
            "options" => {
              "user" => "root"
            }
          }
        }))

        expect(config["Mappings"][0]["Network"]).to eq("my-network")
        expect(config["Mappings"][0]["TLD"]).to eq("my-tld")
    end
  end
end
