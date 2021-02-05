export const currentHash = () => {
  let split = window.location.href.split("/");
  return split.pop() || split.pop();
};
