import {Format} from "@/model/photo";

let webpResult: boolean | undefined = undefined;

function supportsWebp(): Promise<boolean> {
  if (webpResult !== undefined) {
    return Promise.resolve(webpResult);
  } else {
    return new Promise(function(resolve) {
      const img = new Image();
      img.onload = function() {
        const result = img.width > 0 && img.height > 0;
        webpResult = result;
        resolve(result);
      };
      img.onerror = function() {
        webpResult = false;
        resolve(false);
      };
      img.src = "data:image/webp;base64,UklGRiIAAABXRUJQVlA4IBYAAAAwAQCdASoBAAEADsD+JaQAA3AAAAAA";
    });
  }
}

export function formatUrl(image: string, purpose: number, format = "auto", download = false) {
  const downloadQuery = download ? "?download" : "";
  if (format === "auto") {
    return supportsWebp().then(result => {
      if (result) {
        return `/data/${image}/${purpose}.webp${downloadQuery}`;
      } else {
        return `/data/${image}/${purpose}.jpeg${downloadQuery}`;
      }
    });
  } else {
    return Promise.resolve(`/data/${image}/${purpose}.${format}${downloadQuery}`);
  }
}

export function formatDownloadUrl(image: string, format: Format) {
  return formatUrl(image, format.purpose, format.format, true);
}
