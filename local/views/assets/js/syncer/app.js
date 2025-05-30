async function loadFlutterDoctor() {
  const res = await fetch("/api/flutter-doctor");
  const data = await res.json();

  // First update local (so it's not stuck on 'Loading...')
  document.getElementById("doctor-local").innerHTML = formatDoctorOutput(
    data.local
  );

  // Then update server
  document.getElementById("doctor-server").innerHTML = formatDoctorOutput(
    data.server
  );

}

async function loadPackageDiff() {
  const res = await fetch("/api/package-diff");
  const data = await res.json();
  const diffContainer = document.getElementById("package-diff");
  diffContainer.innerHTML = ""; // Clear any previous content

  // Show a message if there are no missing packages
  if (data.diff !== null && data.diff.length === 0) {
    diffContainer.textContent = "No missing packages. Everything is synced.";
  } else {
    if (data.diff_message !== null) {
      const li = document.createElement("li");
      li.textContent = data.diff_message;
      diffContainer.appendChild(li);
    }
    if (data.diff !== null) {
      // Otherwise, show the list of missing packages
      data.diff.forEach((pkg) => {
        const li = document.createElement("li");
        li.textContent = pkg;
        diffContainer.appendChild(li);
      });
    }
  }

  // Optionally update other sections like "Local" or "Server" tabs
  const localPathsContainer = document.getElementById("local-paths");
  const serverPathsContainer = document.getElementById("server-paths");
  localPathsContainer.innerHTML = `<li>${data.local}</li>`;
  serverPathsContainer.innerHTML = `<li>${data.server}</li>`;
}

async function syncPackages() {
  const btn = document.getElementById("sync-btn");
  const timerEl = document.getElementById("sync-timer");

  btn.disabled = true;
  btn.textContent = "Syncing...";

  // Show and start timer
  timerEl.style.display = "inline-block";
  let seconds = 0;
  timerEl.textContent = "Elapsed: 0s";

  const interval = setInterval(() => {
    seconds++;
    timerEl.textContent = `Elapsed: ${seconds}s`;
  }, 1000);

  // Show loading placeholders
  document.getElementById("local-paths").innerHTML = "<li>Waiting for refresh...</li>";
  document.getElementById("server-paths").innerHTML = "<li>Waiting for refresh...</li>";
  document.getElementById("package-diff").innerHTML = "<li>Waiting for refresh...</li>";

  await fetch("/api/upload", { method: "POST" });
  console.log("here 1")
  await sync();
  console.log("here 2")
  await loadPackageDiff();

  clearInterval(interval);
  timerEl.textContent = `Elapsed: ${seconds}s (done)`;
  btn.disabled = false;
  btn.textContent = "Sync Packages";
}


async function sync() {
  const res = await fetch("/api/sync-packages");
  const data = await res.json();
}


async function refreshPackageTabs() {
  // Display "Waiting for refresh..." while data is loading
  document.getElementById("local-paths").innerHTML =
    "<li>Waiting for refresh...</li>";
  document.getElementById("server-paths").innerHTML =
    "<li>Waiting for refresh...</li>";
  document.getElementById("package-diff").innerHTML =
    "<li>Waiting for refresh...</li>";

  try {
    // Fetch the updated package diff data from the server
    const res = await fetch("/api/package-diff");
    const data = await res.json();

    // Handle 'local' field (string or array)
    if (typeof data.local === "string") {
      document.getElementById(
        "local-paths"
      ).innerHTML = `<li>${data.local}</li>`;
    } else if (Array.isArray(data.local)) {
      document.getElementById("local-paths").innerHTML = data.local
        .map((p) => `<li>${p}</li>`)
        .join("");
    }

    // Handle 'server' field (string or array)
    if (typeof data.server === "string") {
      document.getElementById(
        "server-paths"
      ).innerHTML = `<li>${data.server}</li>`;
    } else if (Array.isArray(data.server)) {
      document.getElementById("server-paths").innerHTML = data.server
        .map((p) => `<li>${p}</li>`)
        .join("");
    }

    const diffContainer = document.getElementById("package-diff");
    diffContainer.innerHTML = "";

    // Safeguard: Check if `data.diff` is an array before trying to access its length
    if (
      Array.isArray(data.diff) &&
      data.diff !== null &&
      data.diff.length === 0
    ) {
      diffContainer.innerHTML =
        "<li>No missing packages. Everything is synced.</li>";
    } else {
      if (Array.isArray(data.diff)) {
        if (data.diff_message !== null) {
          const li = document.createElement("li");
          li.textContent = data.diff_message;
          diffContainer.appendChild(li);
        }
        if (data.diff !== null) {
          data.diff.forEach((pkg) => {
            const li = document.createElement("li");
            li.textContent = pkg;
            diffContainer.appendChild(li);
          });
        }
      } else {
        // Handle case when `data.diff` is not available
        diffContainer.innerHTML = "<li>Error: Diff data is not available.</li>";
      }
    }
  } catch (error) {
    console.error("Error fetching data:", error);
    document.getElementById("package-diff").innerHTML =
      "<li>Error: Failed to load package diff.</li>";
  }
}


async function archive() {
  const btn = document.getElementById("archive-btn");
  btn.disabled = true;
  btn.textContent = "Archiving...";

  try {
    const res = await fetch("/api/archive");
    if (!res.ok) {
      throw new Error("Failed to archive packages");
    }
    const data = await res.json(); // { name, size, path }

    document.getElementById("archive-name").textContent = data.name;
    document.getElementById("archive-size").textContent = (data.size / 1024).toFixed(2) + " KB";
    document.getElementById("archive-path").textContent = data.path;

    // Show the archive result
    document.getElementById("archive-result").style.display = "block";

  } catch (err) {
    console.error(err);
    document.getElementById("archive-result").innerHTML = 
      "<p style='color:red;'>Failed to create archive.</p>";
    document.getElementById("archive-result").style.display = "block";
  }

  btn.disabled = false;
  btn.textContent = "Archive Packages";
}
