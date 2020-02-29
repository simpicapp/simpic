import Axios from 'axios'
import { castArray } from 'lodash-es'
import LRUMap from 'mnemonist/lru-map'

class Cache {
  constructor () {
    this._cachedImages = new LRUMap(10)
    this._cachedMetadata = new LRUMap(1000)
    this._cachedThumbnails = new LRUMap(1000)
  }

  storeMetadata (metadata) {
    castArray(metadata).forEach((m) => {
      this._cachedMetadata.set(m.id, Promise.resolve(m))
    })
  }

  getMetadata (id) {
    if (!this._cachedMetadata.has(id)) {
      this._cachedMetadata.set(id, Axios
        .get('/photos/' + id)
        .then(({ data }) => {
          return data
        }))
    }
    return this._cachedMetadata.get(id)
  }

  getThumbnail (id) {
    return this._loadImage(this._cachedThumbnails, '/data/thumb/', id)
  }

  getImage (id) {
    return this._loadImage(this._cachedImages, '/data/image/', id)
  }

  _loadImage (cache, prefix, id) {
    if (!cache.has(id) || cache.peek(id) === null) {
      cache.set(id, new Promise((resolve, reject) => {
        const img = new Image()
        img.onload = () => {
          resolve(img)
        }
        img.onerror = () => {
          cache.set(id, null)
          reject(img)
        }
        img.src = prefix + id
      }))
    }
    return cache.get(id)
  }
}

export const cache = new Cache()
