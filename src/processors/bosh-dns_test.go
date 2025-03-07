package processors

import "testing"

func TestIsAllowedFQDN(t *testing.T) {
	regexps, err := createRegexps([]string{"*.boshdev.lan", "*.appdev.lan"})
	if err != nil {
		t.Error(err)
	}
	processor := boshDnsProcessor{
		fqdnAllowedRegexps: regexps,
	}

	assertTrue := func(fqdn string) {
		if !processor.IsAllowedFQDN(fqdn) {
			t.Errorf("Expected '%s' to match", fqdn)
		}
	}
	assertFalse := func(fqdn string) {
		if processor.IsAllowedFQDN(fqdn) {
			t.Errorf("Expected '%s' to not match", fqdn)
		}
	}

	assertTrue("www.boshdev.lan")
	assertTrue("www.appdev.lan")
	assertFalse("www.boshprod.lan")
	assertFalse("www.appprod.lan")
}
