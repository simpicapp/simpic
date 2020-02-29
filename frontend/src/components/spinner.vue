<template>
  <div class="loader"></div>
</template>

<style scoped>
  .loader,
  .loader:after {
    border-radius: 50%;
    width: 10em;
    height: 10em;
  }

  .loader {
    margin: 60px auto;
    font-size: 10px;
    position: relative;
    text-indent: -9999em;
    border-top: 1.1em solid rgba(74, 153, 118, 0.2);
    border-right: 1.1em solid rgba(74, 153, 118, 0.2);
    border-bottom: 1.1em solid rgba(74, 153, 118, 0.2);
    border-left: 1.1em solid #4a9976;
    transform: translateZ(0);
    animation: load8 1.1s infinite linear;
  }

  @keyframes load8 {
    0% {
      transform: rotate(0deg);
    }
    100% {
      transform: rotate(360deg);
    }
  }
</style>
