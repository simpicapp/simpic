<template>
    <div class="thumbnail" v-bind:class="{ selecting }">
        <a v-bind:href="'/data/photo/' + id" v-on:click.prevent="handleClick">
            <img v-bind:src="'/data/thumb/' + id" v-bind:alt="caption">
            <div class="overlay">
                <p class="caption">{{ caption }}</p>
            </div>
            <span role="button" class="tickbox" v-bind:class="{ selected }" v-on:click.prevent.stop="handleToggle" v-if="$root.loggedIn">
                {{ selected ? '☑' : '☐'}}
            </span>
        </a>
    </div>
</template>

<style>
    .thumbnail {
        position: relative;
        margin: 3px 5px;
    }

    .thumbnail img {
        max-width: 800px;
    }

    .thumbnail:hover .overlay {
        display: grid;
    }

    .thumbnail .overlay {
        display: none;
        position: absolute;
        top: 0;
        bottom: 0;
        left: 0;
        right: 0;
        grid-template-rows: 25% auto 25%;
        grid-template-columns: 25% auto 25%;
        overflow: hidden;
    }

    .thumbnail .tickbox {
        position: absolute;
        top: 0;
        left: 0;
        width: 50px;
        height: 50px;
        display: none;
        background-color: #ffffffdd;
        font-size: xx-large;
        align-items: center;
        justify-content: center;
        border-bottom-right-radius: 2px;
        color: #00000066;
    }

    .thumbnail:hover .tickbox, .thumbnail.selecting .tickbox {
        display: flex;
    }

    .thumbnail .tickbox.selected {
        color: black;
    }

    .thumbnail .caption {
        background-color: #000000cc;
        grid-area: 3 / 1 / 4 / 4;
        text-align: center;
        align-self: end;
        margin: 0;
        padding: 5px;
        color: white;
        font-weight: bold;
        overflow: hidden;
        min-width: 0;
    }
</style>

<script>
  export default {
    props: ['id', 'caption', 'selecting'],
    data () {
      return {
        selected: false
      }
    },
    methods: {
      handleClick () {
        this.$router.push({ path: 'photo/' + this.id })
        this.$emit('showing-photo', this.id)
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
    }
  }
</script>
