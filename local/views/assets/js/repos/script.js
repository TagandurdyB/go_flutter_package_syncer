document.getElementById("cloneForm").addEventListener("submit", async (e) => {
  e.preventDefault();
  const repoUrl = document.getElementById("repoUrl").value;
  const status = document.getElementById("statusMessage");
  status.textContent = "Cloning... Please wait.";

  try {
    const response = await fetch("api/clone", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ repoUrl }),
    });

    const result = await response.json();
    if (response.ok) {
      status.textContent = "✅ Repository cloned successfully!";
    } else {
      status.textContent = "❌ Failed: " + result.error;
    }
  } catch (err) {
    status.textContent = "❌ Error connecting to server.";
  }
});

function openBranches(repoName) {
  fetch(`/api/branches/${repoName}`)
    .then((res) => res.json())
    .then((data) => {
      const card = [...document.querySelectorAll(".repo-card")].find(
        (el) => el.querySelector("h3").textContent === repoName
      );

      let list = document.createElement("ul");
      list.className = "branch-list";

      list.innerHTML = data.branches
        .map(
          (branch) => `
          <li class="branch-item">
            <a class="branch-link" href="/repos/${encodeURIComponent(
              repoName
            )}/${encodeURIComponent(branch)}">
              ${branch}
            </a>
          </li>
        `
        )
        .join("");

      const existing = card.querySelector("ul");
      if (existing) existing.remove();

      card.appendChild(list);
    })
    .catch((err) => {
      alert("Failed to load branches");
      console.error(err);
    });
}


const outputCard = document.getElementById("output");
const outputContent = document.getElementById("output-content");
outputCard.style.display = "block";
const outputTitle = document.getElementById("terminal-header");

function runGitPull(repo, branch) {
  outputContent.textContent = "Pulling from git...";
  outputTitle.textContent = `Pulling from ${repo} (${branch})`;

  fetch(
    `/api/git-pull?repo=${encodeURIComponent(repo)}&branch=${encodeURIComponent(
      branch
    )}`
  )
    .then(async (res) => {
      if (!res.ok) {
        console.log("res:=", res);
        const errorData = await res
          .json()
          .catch(() => ({ message: "Unknown error", error: "Unknown" }));
        throw new Error(
          "❌ Error: " + errorData.error + "\n" + errorData.message
        );
      }
      return res.text();
    })
    .then((data) => {
      outputContent.textContent = data;
    })
    .catch((err) => {
      outputContent.textContent = err.message;
    });
}

function runBuild(repoName, branch, action) {
  const outputCard = document.getElementById("output");
  const outputContent = document.getElementById("output-content");
  outputCard.style.display = "block";
  const outputTitle = document.getElementById("terminal-header");
  outputTitle.textContent = `Running ${action} on ${repoName} (${branch})`;
  outputContent.textContent = "Wating for " + action + "...";

  const endpoint = "/api/build";
  const payload = {
    repo: repoName,
    branch: branch,
    action: action,
  };

  // Disable buttons during request
  const buttons = document.querySelectorAll("button");
  buttons.forEach((btn) => (btn.disabled = true));

  // Optional: show feedback
  console.log(`Running ${action} on ${repoName} [${branch}]`);

  fetch(endpoint, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(payload),
  })
    .then(async (res) => {
      if (!res.ok) {
        console.log("res:=", res);
        const errorData = await res
          .json()
          .catch(() => ({ message: "Unknown error", error: "Unknown" }));
        throw new Error(
          "❌ Error: " + errorData.error + "\n" + errorData.message
        );
      }
      return res.text();
    })
    .then((data) => {
      // outputContent.textContent = data;
      outputContent.innerHTML = formatDoctorOutput(data);
    })
    .catch((err) => {
      outputContent.textContent = err.message;
    });
  // .then((response) => response.json())
  // .then((data) => {
  //   alert(
  //     `✅ ${action} completed!\n\nOutput:\n${data.output || "No output."}`
  //   );
  // })
  // .catch((err) => {
  //   console.error("❌ Build error:", err);
  //   alert(`❌ Error running ${action}: ${err.message}`);
  // })
  // .finally(() => {
  //   buttons.forEach((btn) => (btn.disabled = false));
  // });
}

function flutterPubGet(repoName) {

  outputTitle.textContent = `Running flutter pub get on ${repoName}`;
  outputContent.textContent = "Wating for flutter pub get...";

  const endpoint = `/api/flutter/pub-get${repoName}`;
  // Disable buttons during request
  const buttons = document.querySelectorAll("button");
  buttons.forEach((btn) => (btn.disabled = true));

  fetch(endpoint, {
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ repoName }),
  })
    .then(async (res) => {
      if (!res.ok) {
        console.log("res:=", res);
        const errorData = await res
          .json()
          .catch(() => ({ message: "Unknown error", error: "Unknown" }));
        throw new Error(
          "❌ Error: " + errorData.error + "\n" + errorData.message
        );
      }
      return res.text();
    })
    .then((data) => {
      // outputContent.textContent = data;
      outputContent.innerHTML = formatDoctorOutput(data);
    })
    .catch((err) => {
      outputContent.textContent = err.message;
    });
}
