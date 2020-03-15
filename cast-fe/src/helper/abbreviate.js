export default function abbreviate(value) {
  if (value >= 10 ** 9) {
    value /= 10 ** 9;
    return Math.round(value * 10) / 10 + 'b';
  }
  if (value >= 10 ** 6) {
    value /= 10 ** 6;
    return Math.round(value * 10) / 10 + 'm';
  }
  if (value >= 10 ** 3) {
    value /= 10 ** 3;
    return Math.round(value * 10) / 10 + 'k';
  }
  return value
}
