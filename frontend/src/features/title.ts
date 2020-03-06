const titleParts = ["Simpic"];

export function useTitle() {
  function updateTitle() {
    document.title = titleParts
      .concat()
      .reverse()
      .join(" - ");
  }

  function truncateTitle(max: number) {
    while (titleParts.length > max) {
      titleParts.pop();
    }
    updateTitle();
  }

  function setTitle(index: number, part: string) {
    titleParts[index] = part;
    truncateTitle(index + 1);
  }

  return {truncateTitle, setTitle};
}
