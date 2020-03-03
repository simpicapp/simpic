import LRUMap from "mnemonist/lru-map";

class Cache {
  private readonly _cachedImages: LRUMap<string, Promise<HTMLImageElement>>;
  private readonly _cachedThumbnails: LRUMap<string, Promise<HTMLImageElement>>;

  constructor() {
    this._cachedImages = new LRUMap(10);
    this._cachedThumbnails = new LRUMap(1000);
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
