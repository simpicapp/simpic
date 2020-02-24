<template>
    <div class="thumbnail" :class="{ selecting }" :style="styles">
        <a :href="'/data/image/' + id" @click.prevent="handleClick">
            <div class="overlay">
                <p class="caption">{{ caption }}</p>
            </div>
            <span role="button" class="tickbox"
                  :class="{ selected }"
                  @click.prevent.stop="handleToggle"
                  v-if="$root.loggedIn">
                {{ selected ? '☑' : '☐'}}
            </span>
        </a>
    </div>
</template>

<style lang="scss" scoped>
    .thumbnail {
        flex-grow: 1;
        flex-shrink: 1;
        position: relative;
        height: 200px;
        margin: 3px 3px;
        background-repeat: no-repeat;
        background-position: 50%;
        background-size: cover;
    }

    a {
        display: block;
        height: 200px;
    }

    img {
        max-width: 800px;
    }

    .overlay {
        display: grid;
        position: absolute;
        top: 0;
        bottom: 0;
        left: 0;
        right: 0;
        grid-template-rows: 25% auto 25%;
        grid-template-columns: 25% auto 25%;
        overflow: hidden;
        opacity: 0;
        transition: opacity 300ms var(--ease-in-cubic);
    }

    .tickbox {
        position: absolute;
        top: 0;
        left: 0;
        width: 50px;
        height: 50px;
        display: flex;
        background-color: #ffffffdd;
        font-size: xx-large;
        align-items: center;
        justify-content: center;
        border-bottom-right-radius: 5px;
        opacity: 0;
        transition: opacity 300ms var(--ease-in-cubic);

        color: #00000066;

        &.selected {
            color: black;
        }
    }

    .thumbnail:hover .overlay, .thumbnail:hover .tickbox, .thumbnail.selecting .tickbox {
        opacity: 1;
        transition: opacity 300ms var(--ease-out-cubic);
    }

    .caption {
        grid-area: 3 / 1 / 4 / 4;
        text-align: center;
        align-self: end;
        margin: 0;
        padding: 10px 0;
        color: white;
        overflow: hidden;
        min-width: 0;

        @supports (backdrop-filter: blur()) {
            backdrop-filter: blur(10px);
            background-color: #00000099;
        }

        @supports not (backdrop-filter: blur()) {
            background-color: #000000cc;
        }
    }
</style>

<script>
  export default {
    props: ['id', 'caption', 'selecting', 'selected'],
    data () {
      return {
        styles: {
          backgroundImage: '',
          flexBasis: 0,
          maxWidth: 0
        }
      }
    },
    methods: {
      handleClick (e) {
        if (this.selecting && e.ctrlKey) {
          // Ctrl+click during selection is a shortcut for toggling
          this.handleToggle()
        } else if (this.selecting && e.shiftKey) {
          // Shift+click is a shortcut for range selection
          this.$emit('select-range', this.id)
        } else {
          this.$router.push({ path: 'photo/' + this.id })
          this.$emit('showing-photo', this.id)
        }
      },
      handleImageLoaded (e) {
        console.log(e)
      },
      handleToggle () {
        if (this.selected) {
          this.$emit('deselected', this.id)
        } else {
          this.$emit('selected', this.id)
        }
      }
    },
    mounted () {
      const image = new Image()
      image.onload = () => {
        const canvas = document.createElement('canvas')
        canvas.width = image.naturalWidth
        canvas.height = image.naturalHeight
        canvas.getContext('2d').drawImage(image, 0, 0)

        this.styles.maxWidth = (image.naturalWidth * 1.5) + 'px'
        this.styles.flexBasis = image.naturalWidth + 'px'
        this.styles.backgroundImage = 'url(' + canvas.toDataURL('image/jpeg') + ')'
      }
      image.src = '/data/thumb/' + this.id
    }
  }
</script>
