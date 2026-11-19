export function toast(msg) {
		const p = document.createElement("p");
		p.textContent = String(msg).charAt(0).toUpperCase() + String(msg).slice(1);
		p.classList.add("toast");
		document.body.appendChild(p);
		setTimeout(() => {
				p.remove();
		}, 2400);
}