export function getInputWidth(value: string, font: string = "16px Arial") {
  // create virtual span
  const input = document.createElement("input");
  input.style.visibility = "hidden";
  input.style.position = "absolute";
  input.style.font = font;
  input.value = value || " ";
  document.body.appendChild(input);

  const width = input.offsetWidth + 8;
  document.body.removeChild(input);

  return Math.max(width, 32);
}
