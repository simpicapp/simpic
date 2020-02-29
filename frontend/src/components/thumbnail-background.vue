<script lang="ts">
  import {cache} from './cache'
  import Vue from 'vue'

  export default Vue.extend({
    props: ['imageId'],
    data() {
      return {
        styles: {
          backgroundImage: '',
          flexBasis: '0',
          maxWidth: '0'
        }
      }
    },
    methods: {
      showFallbackStyle() {
        const angle = Math.random() * 360;
        this.styles = {
          backgroundImage: `repeating-linear-gradient(${angle}deg, #244c3b, #244c3b 10px, #183327 10px, #183327 20px)`,
          flexBasis: '150px',
          maxWidth: '300px'
        }
      }
    },
    mounted() {
      if (this.imageId) {
        cache.getThumbnail(this.imageId).then((image: HTMLImageElement) => {
          const canvas = document.createElement('canvas');
          const ctx = canvas.getContext('2d');
          canvas.width = image.naturalWidth;
          canvas.height = image.naturalHeight;
          ctx && ctx.drawImage(image, 0, 0);

          this.styles = {
            backgroundImage: 'url(' + canvas.toDataURL('image/jpeg') + ')',
            flexBasis: image.naturalWidth + 'px',
            maxWidth: (image.naturalWidth * 1.5) + 'px'
          }
        }).catch(() => this.showFallbackStyle())
      } else {
        this.showFallbackStyle()
      }
    }
  })
</script>
