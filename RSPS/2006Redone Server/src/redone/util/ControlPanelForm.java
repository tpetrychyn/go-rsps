package redone.util;

import lombok.Getter;
import redone.Server;
import redone.game.objects.Object;
import redone.game.players.Player;
import redone.game.players.PlayerHandler;

import javax.swing.*;
import java.util.Arrays;

@Getter
public class ControlPanelForm {
    private Player player;
    private JTabbedPane tabbedPane1;
    private JPanel mainPanel;
    private JTextField objectId;
    private JButton spawnButton;
    private JToolBar playerInfo;
    private JLabel playerCoords;
    private JButton searchButton;
    private JTextField npcId;
    private JButton spawnNpc;
    private JList coordsList;
    private JButton addToCoordsList;

    public void setPlayer(Player p) {
        this.player = p;
    }

    ControlPanelForm() {
        searchButton.addActionListener(e ->
                Arrays.asList(PlayerHandler.players).forEach(p -> {
                    if (p != null) {
                        setPlayer(p);
                    }
                }));
        spawnButton.addActionListener(e -> {
            int id = Integer.parseInt(objectId.getText());
            new Object(id, player.absX, player.absY, player.heightLevel, player.face, 10, -1, 100);
        });
        spawnNpc.addActionListener(e -> {
            int id = Integer.parseInt(npcId.getText());
            Server.npcHandler.spawnNewNpc(id, player.absX, player.absY, player.heightLevel, 0, 20, 2, 0, 1000);
        });
        addToCoordsList.addActionListener(e -> {

        });

    }

    public void process() {
        if (player == null) return;

        playerCoords.setText("Pos: " + player.absX + ", " + player.absY);
    }
}
