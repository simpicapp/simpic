import {onMounted, onUnmounted} from "@vue/composition-api";
import {throttle} from "lodash-es";
import {EventBus} from "@/components/bus";

export function usesScrollWatcher() {
  let atBottom = false;

  const emitBottom = throttle(() => EventBus.$emit('bottom'), 250);

  function bottomVisible() {
    const scrollY = window.scrollY;
    const visible = document.documentElement.clientHeight;
    const pageHeight = document.documentElement.scrollHeight;
    return visible + scrollY >= pageHeight - 400
  }

  function handleScroll() {
    const nowAtBotton = bottomVisible();
    if (nowAtBotton && !atBottom) {
      emitBottom()
    }
    atBottom = nowAtBotton;
  }

  onUnmounted(() => {
    window.removeEventListener('scroll', handleScroll)
  });

  onMounted(() => {
    window.addEventListener('scroll', handleScroll)
  });
}
