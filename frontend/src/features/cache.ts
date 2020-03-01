import Axios from "axios";
import {castArray} from "lodash-es";
import LRUMap from "mnemonist/lru-map";
import {Photo} from "@/model/photo";

class Cache {
  private readonly _cachedImages: LRUMap<string, Promise<HTMLImageElement>>;
  private readonly _cachedThumbnails: LRUMap<string, Promise<HTMLImageElement>>;
  private readonly _cachedMetadata: LRUMap<string, Promise<Photo>>;

  constructor() {
    this._cachedImages = new LRUMap(10);
    this._cachedMetadata = new LRUMap(1000);
    this._cachedThumbnails = new LRUMap(1000);
  }

  storeMetadata(metadata: Photo | Array<Photo>) {
    castArray(metadata).forEach(m => {
      this._cachedMetadata.set(m.id, Promise.resolve(m));
    });
  }

  getMetadata(id: string): Promise<Photo> {
    const cached = this._cachedMetadata.get(id);
    if (cached) {
      return cached;
    }

    const created = Axios.get("/photos/" + id).then(({data}) => {
      return data;
    });

    this._cachedMetadata.set(id, created);
    return created;
  }

  getThumbnail(id: string) {
    return this._loadImage(this._cachedThumbnails, "/data/thumb/", id);
  }

  getImage(id: string) {
    return this._loadImage(this._cachedImages, "/data/image/", id);
  }

  _loadImage(
    cache: LRUMap<string, Promise<HTMLImageElement> | null>,
    prefix: string,
    id: string
  ): Promise<HTMLImageElement> {
    const cached = cache.get(id);
    if (cached) {
      return cached;
    }

    const created = new Promise<HTMLImageElement>((resolve, reject) => {
      const img = new Image();
      img.onload = () => {
        resolve(img);
      };
      img.onerror = () => {
        cache.set(id, null);
        reject(img);
      };
      img.src = prefix + id;
    });
    cache.set(id, created);
    return created;
  }
}

export const cache = new Cache();
