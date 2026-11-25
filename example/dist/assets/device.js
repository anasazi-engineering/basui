async function loadDevice() {
  const id = new URLSearchParams(location.search).get("id");
  const status = document.querySelector("#status");
  if (!id) {
    status.textContent = "Select a device from the equipment page.";
    return;
  }

  try {
    const response = await fetch(new URL("./devices.json", import.meta.url));
    if (!response.ok) throw new Error("Could not load readings.");
    const devices = await response.json();
    const device = devices.find(item => item.id === id);
    if (!device) throw new Error("Device not found.");

    document.querySelector("#name").textContent = device.name;
    const template = document.querySelector("#point-row");
    for (const point of device.points) {
      const row = template.content.firstElementChild.cloneNode(true);
      row.cells[0].textContent = point.name;
      row.cells[1].textContent = point.value ?? "—";
      row.cells[2].textContent = point.units ?? "";
      document.querySelector("#points").append(row);
    }
    status.textContent = "Sample readings";
  } catch (error) {
    status.textContent = error.message;
  }
}

loadDevice();
