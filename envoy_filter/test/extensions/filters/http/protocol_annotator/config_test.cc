#include "extensions/filters/http/protocol_annotator/config.h"

#include "test/mocks/server/mocks.h"

namespace Envoy {
namespace Extensions {
namespace HttpFilters {
namespace ProtocolAnnotator {

TEST(ProtocolAnnotatorConfigFactoryTest, Config) {
  const std::string yaml = R"EOF(
      header_name: x-envoy-annotated-protocol
    )EOF";

  io::omnition::envoy::v1::ProtocolAnnotator proto_config;
  MessageUtil::loadFromYaml(yaml, proto_config);

  NiceMock<Server::Configuration::MockFactoryContext> context;
  ProtocolAnnotatorFilterConfigFactory factory;

  Http::FilterFactoryCb cb = factory.createFilterFactoryFromProto(proto_config, "stats", context);
  Http::MockFilterChainFactoryCallbacks filter_callback;
  EXPECT_CALL(filter_callback, addStreamDecoderFilter(testing::_));
  cb(filter_callback);
}

} // namespace ProtocolAnnotator
} // namespace HttpFilters
} // namespace Extensions
} // namespace Envoy
