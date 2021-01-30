const month = [
  "January",
  "February",
  "March",
  "April",
  "May",
  "June",
  "July",
  "August",
  "September",
  "October",
  "November",
  "December",
];

export default function format() {
  return {
    date: (dateString) => {
      let date = new Date(dateString);
      return `${date.getDate()} ${
        month[date.getMonth()]
      } ${date.getFullYear()}`;
    },
    full_date: (dateString) => {
      let date = new Date(dateString);
      return `${date.getHours()}:${
        date.getMinutes() < 10 ? "0" : ""
      }${date.getMinutes()}, ${date.getDate()} ${
        month[date.getMonth()]
      } ${date.getFullYear()}`;
    },
  };
}
