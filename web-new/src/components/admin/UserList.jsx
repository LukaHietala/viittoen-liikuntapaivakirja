import UserItem from "./UserItem.jsx";

function UserList(props) {
    return (
        <div>
            <For each={props.users()}>{(item) =>
                <UserItem
										user={item}
										editUser={props.editUser}
										deleteUser={props.deleteUser}
								/>
            }</For>
        </div>
    );
}

export default UserList;