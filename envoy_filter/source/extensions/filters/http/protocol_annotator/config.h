#pragma once

#include "source/extensions/filters/http/protocol_annotator/config.pb.validate.h"

#include "extensions/filters/http/common/factory_base.h"

namespace Envoy {
namespace Extensions {
namespace HttpFilters {
namespace ProtocolAnnotator {

static const std::string& PROTOCOL_ANNOTATOR_FILTER() {
  // TODO(dio): Should move this to a dedicated well_known_names.h header file.
  CONSTRUCT_ON_FIRST_USE(std::string, "io.omnition.envoy.http.protocol_annotator");
}

class ProtocolAnnotatorFilterConfigFactory
    : public Common::FactoryBase<io::omnition::envoy::v1::ProtocolAnnotator> {
public:
  ProtocolAnnotatorFilterConfigFactory() : FactoryBase(PROTOCOL_ANNOTATOR_FILTER()) {}

private:
  Http::FilterFactoryCb
  createFilterFactoryFromProtoTyped(const io::omnition::envoy::v1::ProtocolAnnotator& proto_config,
                                    const std::string& stats_prefix,
                                    Server::Configuration::FactoryContext& context) override;
};

} // namespace ProtocolAnnotator
} // namespace HttpFilters
} // namespace Extensions
} // namespace Envoy
