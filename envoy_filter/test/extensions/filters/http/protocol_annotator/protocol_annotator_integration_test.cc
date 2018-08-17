#include "test/integration/http_integration.h"
#include "test/integration/utility.h"

namespace Envoy {
class ProtocolAnnotatorFilterIntegrationTest
    : public HttpIntegrationTest,
      public testing::TestWithParam<Network::Address::IpVersion> {

public:
  ProtocolAnnotatorFilterIntegrationTest()
      : HttpIntegrationTest(Http::CodecClient::Type::HTTP1, GetParam()) {}

  void SetUp() override { initialize(); }

  void initialize() override {
    config_helper_.addFilter(R"EOF(
    name: io.omnition.envoy.http.protocol_annotator
    config:
      header_name: x-envoy-annotated-protocol
    )EOF");

    HttpIntegrationTest::initialize();
  }

protected:
  void testRequest(Http::TestHeaderMapImpl&& request_headers,
                   Http::TestHeaderMapImpl&& expected_response_headers,
                   bool wait_for_upstream_response = false) {
    codec_client_ = makeHttpConnection(lookupPort("http"));
    if (wait_for_upstream_response) {
      sendRequestAndWaitForResponse(request_headers, 0, expected_response_headers, 0);
      EXPECT_EQ("HTTP/1.1", static_cast<Http::TestHeaderMapImpl>(upstream_request_->headers())
                                .get_("x-envoy-annotated-protocol"));
    } else {
      auto response = codec_client_->makeHeaderOnlyRequest(request_headers);
      // and don't wait for upstream response
      response->waitForEndStream();
      EXPECT_TRUE(response->complete());
      compareHeaders(response->headers(), expected_response_headers);
    }
  }

  void compareHeaders(Http::TestHeaderMapImpl&& response_headers,
                      Http::TestHeaderMapImpl& expected_response_headers) {
    response_headers.remove(Envoy::Http::LowerCaseString{"date"});
    response_headers.remove(Envoy::Http::LowerCaseString{"x-envoy-upstream-service-time"});
    EXPECT_EQ(expected_response_headers, response_headers);
  }
};

INSTANTIATE_TEST_CASE_P(IpVersions, ProtocolAnnotatorFilterIntegrationTest,
                        testing::ValuesIn(TestEnvironment::getIpVersionsForTest()));

TEST_P(ProtocolAnnotatorFilterIntegrationTest, Passthrough) {
  testRequest(
      Http::TestHeaderMapImpl{{":method", "GET"}, {":path", "/path"}, {":authority", "host"}},
      Http::TestHeaderMapImpl{{"content-length", "0"}, {"server", "envoy"}, {":status", "200"}},
      true);
}

} // namespace Envoy
