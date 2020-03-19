export default function abbreviate() {
  return {
    number: (value) => {
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
    },
    time: (diff) => {
      diff = Math.round(diff / 1000);
      if (diff > 29030400) {
        diff = Math.round(diff / 29030400);
        return `${diff} year${diff === 1 ? "" : "s"}`
      }
      if (diff > 2419200) {
        diff = Math.round(diff / 2419200);
        return `${diff} month${diff === 1 ? "" : "s"}`
      }
      if (diff > 604600) {
        diff = Math.round(diff / 604600);
        return `${diff} week${diff === 1 ? "" : "s"}`
      }
      if (diff > 86400) {
        diff = Math.round(diff / 86400);
        return diff === 1 ? "Yesterday" : `${diff} days`
      }
      if (diff > 3600) {
        diff = Math.round(diff / 3600);
        return `${diff} hour${diff === 1 ? "" : "s"}`
      }
      if (diff > 60) {
        diff = Math.round(diff / 60);
        return `${diff} minute${diff === 1 ? "" : "s"}`
      }
      return `${diff} second${diff === 1 ? "" : "s"}`
    }
  }
}
