import {Format, PurposeDownload, PurposePreview, PurposeScreen} from "@/model/photo";

export function formatFileSize(bytes: number) {
  const thresh = 1024;
  if (Math.abs(bytes) < thresh) {
    return bytes + " B";
  }
  const units = ["KiB", "MiB", "GiB", "TiB", "PiB", "EiB", "ZiB", "YiB"];
  let u = -1;
  do {
    bytes /= thresh;
    ++u;
  } while (Math.abs(bytes) >= thresh && u < units.length - 1);
  return bytes.toFixed(1) + " " + units[u];
}

export function formatPurpose(purpose: number) {
  switch (purpose) {
    case PurposeDownload:
      return "Original";
    case PurposeScreen:
      return "Screen-optimised";
    case PurposePreview:
      return "Thumbnail";
    default:
      return "Unknown";
  }
}

export function formatDownloadUrl(image: string, format: Format) {
  let url = "";
  switch (format.purpose) {
    case PurposeScreen:
      url = "image";
      break;
    case PurposeDownload:
      url = "raw";
      break;
    case PurposePreview:
      url = "thumb";
      break;
  }
  return `/data/${url}/${image}?download`;
}
