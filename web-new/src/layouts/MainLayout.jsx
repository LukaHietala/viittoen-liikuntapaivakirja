import { useMatch } from "@solidjs/router";
import { useAuth } from "../context/AuthContext";

function Navbar() {
		const { user } = useAuth();
		
		return (
				<nav class="navbar navbar-row">
						<div>
								<NavLink href="/">Etusivu</NavLink>
								<NavLink href="/profile">Profiili</NavLink>
                <NavLink href="/challenges">Haasteet</NavLink>
								<Show when={user() && user().isAdmin}>
										<NavLink href="/admin">Ylläpito</NavLink>
								</Show>
						</div>
						<div>
								<NavLink href="/logout">Kirjaudu ulos</NavLink>
						</div>
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
		return (
				<>
						<main class="container">
								<Show when={user()}>{() => <Navbar />}</Show>
								{props.children}
						</main>
				</>
		);
}

export default MainLayout;
