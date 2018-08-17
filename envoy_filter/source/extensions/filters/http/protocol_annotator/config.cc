#include "envoy/registry/registry.h"

#include "extensions/filters/http/protocol_annotator/config.h"
#include "extensions/filters/http/protocol_annotator/protocol_annotator.h"

namespace Envoy {
namespace Extensions {
namespace HttpFilters {
namespace ProtocolAnnotator {

Http::FilterFactoryCb ProtocolAnnotatorFilterConfigFactory::createFilterFactoryFromProtoTyped(
    const io::omnition::envoy::v1::ProtocolAnnotator& proto_config, const std::string&,
    Server::Configuration::FactoryContext&) {
  ProtocolAnnotatorFilterConfigPtr config =
      std::make_shared<ProtocolAnnotatorFilterConfig>(ProtocolAnnotatorFilterConfig(proto_config));

  return [config](Http::FilterChainFactoryCallbacks& callbacks) -> void {
    callbacks.addStreamDecoderFilter(
        Http::StreamDecoderFilterSharedPtr{new ProtocolAnnotatorFilter(config)});
  };
}

static Registry::RegisterFactory<ProtocolAnnotatorFilterConfigFactory,
                                 Server::Configuration::NamedHttpFilterConfigFactory>
    registered_;

} // namespace ProtocolAnnotator
} // namespace HttpFilters
} // namespace Extensions
} // namespace Envoy
