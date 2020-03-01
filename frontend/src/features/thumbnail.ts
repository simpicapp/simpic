import {onMounted, ref} from "@vue/composition-api";
import {cache} from "@/features/cache";

export function useThumbnail(id: string | undefined) {
  const styles = ref({});

  function showFallbackStyle() {
    const angle = Math.random() * 360;
    styles.value = {
      backgroundImage: `repeating-linear-gradient(${angle}deg, #244c3b, #244c3b 10px, #183327 10px, #183327 20px)`,
      flexBasis: "150px",
      maxWidth: "300px",
    };
  }

  onMounted(() => {
    if (id) {
      cache
        .getThumbnail(id)
        .then((image: HTMLImageElement) => {
          const canvas = document.createElement("canvas");
          const ctx = canvas.getContext("2d");
          canvas.width = image.naturalWidth;
          canvas.height = image.naturalHeight;
          ctx && ctx.drawImage(image, 0, 0);

          styles.value = {
            backgroundImage: "url(" + canvas.toDataURL("image/jpeg") + ")",
            flexBasis: image.naturalWidth + "px",
            maxWidth: image.naturalWidth * 1.5 + "px",
          };
        })
        .catch(() => showFallbackStyle());
    } else {
      showFallbackStyle();
    }
  });

  return {styles};
}
