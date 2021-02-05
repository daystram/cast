export default function urls() {
  return {
    current_hash: () => {
      let split = window.location.href.split("/");
      return split.pop() || split.pop();
    },
  };
}
