package redone.util;

import javax.swing.*;

public class ControlPanel {
    private ControlPanelForm mainForm;

    public ControlPanel() {
        JFrame frame = new JFrame("My gggg");
        frame.setDefaultCloseOperation(JFrame.EXIT_ON_CLOSE);
        frame.setSize(500, 450);
        frame.setLocation(1430, 0);
        mainForm = new ControlPanelForm();
        frame.add(mainForm.getMainPanel());
        frame.setVisible(true);
    }

    public void process() {
        if (mainForm == null) return;

        mainForm.process();
    }
}
