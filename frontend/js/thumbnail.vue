<template>
    <div class="thumbnail" v-bind:class="{ selecting }" v-bind:style="styles">
        <a v-bind:href="'/data/image/' + id" v-on:click.prevent="handleClick">
            <div class="overlay">
                <p class="caption">{{ caption }}</p>
            </div>
            <span role="button" class="tickbox"
                  v-bind:class="{ selected }"
                  v-on:click.prevent.stop="handleToggle"
                  v-if="$root.loggedIn">
                {{ selected ? '☑' : '☐'}}
            </span>
        </a>
    </div>
</template>

<style scoped>
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

    .thumbnail img {
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
        color: #00000066;
        opacity: 0;
        transition: opacity 300ms var(--ease-in-cubic);
    }

    .thumbnail:hover .overlay, .thumbnail:hover .tickbox, .thumbnail.selecting .tickbox  {
        opacity: 1;
        transition: opacity 300ms var(--ease-out-cubic);
    }

    /*noinspection CssUnusedSymbol*/
    .tickbox.selected {
        color: black;
    }

    .caption {
        backdrop-filter: blur(10px);
        background-color: #00000099;
        grid-area: 3 / 1 / 4 / 4;
        text-align: center;
        align-self: end;
        margin: 0;
        padding: 10px 0;
        color: white;
        overflow: hidden;
        min-width: 0;
    }

    @supports not (backdrop-filter: blur()) {
        .caption {
            background-color: #000000cc;
        }
    }
</style>

<script>
  export default {
    props: ['id', 'caption', 'selecting'],
    data () {
      return {
        selected: false,
        styles: {
          backgroundImage: '',
          flexBasis: 0,
          maxWidth: 0
        }
      }
    },
    methods: {
      handleClick () {
        this.$router.push({ path: 'photo/' + this.id })
        this.$emit('showing-photo', this.id)
      },
      handleImageLoaded (e) {
        console.log(e)
      },
      handleToggle () {
        this.selected = !this.selected
        if (this.selected) {
          this.$emit('selected', this.id)
        } else {
          this.$emit('deselected', this.id)
        }
      }
    },
    watch: {
      selecting (newVal) {
        if (!newVal) {
          this.selected = false
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
