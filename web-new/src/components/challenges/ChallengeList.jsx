import ChallengeItem from "./ChallengeItem.jsx";

function ChallengeList(props) {
    return (
        <>
            <For each={props.challenges()}>{(item) =>
                <ChallengeItem
                    challenge={item}
                    editChallenge={props.editChallenge}
                    deleteChallenge={props.deleteChallenge}
                    deletePerformance={props.deletePerformance}
										users={props.users}
                />
            }</For>
            <Show when={props.challenges().length == 0}>
                <div class="hierarchical-item">
                    Ei haasteita
                </div>
            </Show>
        </>
    );
}

export default ChallengeList;