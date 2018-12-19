#include "extensions/filters/http/protocol_annotator/protocol_annotator.h"

#include "absl/strings/match.h"
#include "common/http/utility.h"

namespace Envoy {
namespace Extensions {
namespace HttpFilters {
namespace ProtocolAnnotator {

Http::FilterHeadersStatus ProtocolAnnotatorFilter::decodeHeaders(Http::HeaderMap& headers, bool) {
  const auto& protocol = decoder_callbacks_->requestInfo().protocol();
  const std::string header_value =
      protocol.has_value() ? Http::Utility::getProtocolString(protocol.value()) : "-";
  headers.addCopy(config_->header_name_, header_value);
  return Http::FilterHeadersStatus::Continue;
}

} // namespace ProtocolAnnotator
} // namespace HttpFilters
} // namespace Extensions
} // namespace Envoy
