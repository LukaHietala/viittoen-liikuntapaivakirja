import { useMatch } from "@solidjs/router";
import { useAuth } from "../context/AuthContext";
import { createSignal } from "solid-js";
import innostunutKamu from "../assets/innostunut-kamu.webp"
import okKamu from "../assets/ok-kamu.webp"
import menuIcon from "../assets/menu.svg"

function Navbar() {
    const { user } = useAuth();
    const [showNav, setShowNav] = createSignal(false);
    
    return (
				<nav class="navbar">
						<div class="navbar-row">
								<div class="desktop">
										<NavLink href="/">Etusivu</NavLink>
										<NavLink href="/profile">Profiili</NavLink>
										<NavLink href="/challenges">Haasteet</NavLink>
										<Show when={user() && user().isAdmin}>
												<NavLink href="/admin">Ylläpito</NavLink>
										</Show>
								</div>
								<div class="desktop">
										<NavLink href="/logout">Kirjaudu ulos</NavLink>
								</div>
								<div class="mobile"></div>
								<div class="mobile">
										<button onClick={() =>setShowNav(!showNav())} class="btn-icon">
												<img src={menuIcon} alt="Menu"  />           
										</button>
								</div>
						</div>
						<Show when={showNav()}>
								<div class="navbar-col"
										 onClick={() => showNav() ? setShowNav(false) : null}>
										<NavLink href="/">Etusivu</NavLink>
										<NavLink href="/profile">Profiili</NavLink>
										<NavLink href="/challenges">Haasteet</NavLink>
										<Show when={user() && user().isAdmin}>
												<NavLink href="/admin">Ylläpito</NavLink>
										</Show>
										<NavLink href="/logout">Kirjaudu ulos</NavLink>
								</div>
						</Show>
				</nav>
    );
}

function NavLink(props) {
    const match = useMatch(() => props.href);

    return (
        <a
            class="btn-link"
            href={props.href}
            classList={{ "active-nav-item": Boolean(match()) }}
        >
            {props.children}
        </a>
    );
}

function MainLayout(props) {
    const { user } = useAuth();
    
    const kamut = [
        innostunutKamu, okKamu
    ];

    return (
        <>
            <main class="container">
                <Show when={user()}>{() => <Navbar />}</Show>
                {props.children}
            </main>
            <img class="kamu" src={kamut[Math.floor(Math.random() * kamut.length)]} alt="Tyytyväinen kamu" />
        </>
    );
}

export default MainLayout;
