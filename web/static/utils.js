function toast(msg) {
		const p = document.createElement("p");
		p.textContent = msg;
		p.classList.add("toast");
		document.body.appendChild(p);
		setTimeout(()=>{p.remove();}, 2400);
}