<template>
    <footer>
        <a href="https://simpic.app/">Simpic - simple self-hosted picture manager</a> &middot;
        <span class="version">{{ versionName }}</span>
        <a :href="`https://github.com/simpicapp/simpic/commit/${$root.gitSHA}`" class="sha" v-if="$root.gitSHA">({{ $root.gitSHA }})</a>
    </footer>
</template>

<style lang="scss" scoped>
    footer {
        border-top: 1px solid #ccc;
        background-color: #eee;
        margin-top: 10px;
        height: 70px;
        display: flex;
        align-items: center;
        column-gap: 10px;
        padding-left: 20px;
    }

    a {
        color: black;
        transition: color 300ms linear;

        &:hover {
            color: blue;
        }

        &.sha {
            color: #999999;
            text-decoration: none;

            &:hover {
                color: black;
            }
        }
    }

    .version {
        color: #333;
    }
</style>

<script>
  export default {
    computed: {
      versionName () {
        if (this.$root.gitTag === '') {
          return 'Unknown release'
        } else if (this.$root.gitTag === this.$root.gitSHA) {
          return 'Pre-release build'
        } else {
          const parts = this.$root.gitTag.split('-')
          return parts[0] + (parts.length > 1 ? '+changes' : '')
        }
      }
    }
  }
</script>
