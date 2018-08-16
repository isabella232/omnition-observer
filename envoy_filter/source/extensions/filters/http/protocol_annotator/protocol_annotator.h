#pragma once

#include "envoy/http/filter.h"

#include "common/common/base64.h"

#include "source/extensions/filters/http/protocol_annotator/config.pb.h"

namespace Envoy {
namespace Extensions {
namespace HttpFilters {
namespace ProtocolAnnotator {

struct ProtocolAnnotatorFilterConfig {
  ProtocolAnnotatorFilterConfig(const io::omnition::envoy::v1::ProtocolAnnotator& proto)
      : header_name_(proto.header_name()) {}

  const Http::LowerCaseString header_name_;
};
typedef std::shared_ptr<ProtocolAnnotatorFilterConfig> ProtocolAnnotatorFilterConfigPtr;

class ProtocolAnnotatorFilter : public Http::StreamDecoderFilter {
public:
  ProtocolAnnotatorFilter(ProtocolAnnotatorFilterConfigPtr config) : config_(config) {}
  ~ProtocolAnnotatorFilter() {}

  // Http::StreamFilterBase
  void onDestroy() override {}

  // Http::StreamDecoderFilter
  Http::FilterHeadersStatus decodeHeaders(Http::HeaderMap&, bool) override;
  Http::FilterDataStatus decodeData(Buffer::Instance&, bool) override {
    return Http::FilterDataStatus::Continue;
  }
  Http::FilterTrailersStatus decodeTrailers(Http::HeaderMap&) override {
    return Http::FilterTrailersStatus::Continue;
  }
  void setDecoderFilterCallbacks(Http::StreamDecoderFilterCallbacks& callback) override {
    decoder_callbacks_ = &callback;
  }

private:
  Http::StreamDecoderFilterCallbacks* decoder_callbacks_;
  ProtocolAnnotatorFilterConfigPtr config_;
};

} // namespace ProtocolAnnotator
} // namespace HttpFilters
} // namespace Extensions
} // namespace Envoy