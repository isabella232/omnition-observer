#include "extensions/filters/http/protocol_annotator/protocol_annotator.h"

#include "test/mocks/http/mocks.h"

namespace Envoy {
namespace Extensions {
namespace HttpFilters {
namespace ProtocolAnnotator {

class ProtocolAnnotatorFilterTest : public testing::Test {
protected:
  void SetUp() override {
    const std::string yaml = R"EOF(
      header_name: x-envoy-annotated-protocol
    )EOF";

    io::omnition::envoy::v1::ProtocolAnnotator proto_config;
    MessageUtil::loadFromYaml(yaml, proto_config);

    config_.reset(new ProtocolAnnotatorFilterConfig(proto_config));
    filter_.reset(new ProtocolAnnotatorFilter(config_));
    filter_->setDecoderFilterCallbacks(callbacks_);
  }

  testing::NiceMock<Http::MockStreamDecoderFilterCallbacks> callbacks_;
  std::unique_ptr<ProtocolAnnotatorFilter> filter_;
  ProtocolAnnotatorFilterConfigPtr config_;
  NiceMock<Envoy::RequestInfo::MockRequestInfo> request_info_;
};

TEST_F(ProtocolAnnotatorFilterTest, AnnotatedRequest) {
  Http::TestHeaderMapImpl request_headers{{":method", "get"}};
  EXPECT_CALL(request_info_, protocol()).WillRepeatedly(testing::Return(Http::Protocol::Http11));
  EXPECT_CALL(callbacks_, requestInfo()).WillOnce(testing::ReturnRef(request_info_));
  EXPECT_EQ(Http::FilterHeadersStatus::Continue, filter_->decodeHeaders(request_headers, false));
  EXPECT_EQ("HTTP/1.1", request_headers.get_("x-envoy-annotated-protocol"));
}

} // namespace ProtocolAnnotator
} // namespace HttpFilters
} // namespace Extensions
} // namespace Envoy
