import LRUMap from "mnemonist/lru-map";
import {PurposePreview, PurposeScreen} from "@/model/photo";
import {formatUrl} from "@/features/images";

class Cache {
  private readonly _cachedImages: LRUMap<string, Promise<HTMLImageElement>>;
  private readonly _cachedThumbnails: LRUMap<string, Promise<HTMLImageElement>>;

  constructor() {
    this._cachedImages = new LRUMap(10);
    this._cachedThumbnails = new LRUMap(1000);
  }

  getThumbnail(id: string) {
    return this._loadImage(this._cachedThumbnails, PurposePreview, id);
  }

  getImage(id: string) {
    return this._loadImage(this._cachedImages, PurposeScreen, id);
  }

  _loadImage(
    cache: LRUMap<string, Promise<HTMLImageElement> | null>,
    purpose: number,
    id: string
  ): Promise<HTMLImageElement> {
    const cached = cache.get(id);
    if (cached) {
      return cached;
    }

    const created = formatUrl(id, purpose).then(url => {
      return new Promise<HTMLImageElement>((resolve, reject) => {
        const img = new Image();
        img.onload = () => {
          resolve(img);
        };
        img.onerror = () => {
          cache.set(id, null);
          reject(img);
        };
        img.src = url;
      });
    });
    cache.set(id, created);
    return created;
  }
}

export const cache = new Cache();
