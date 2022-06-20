package com.example.summer_ctf;

import androidx.appcompat.app.AppCompatActivity;

import android.nfc.Tag;
import android.os.Bundle;
import android.util.Log;
import android.view.View;
import android.widget.Button;
import android.widget.EditText;
import android.widget.Toast;

import java.math.BigInteger;
import java.nio.charset.StandardCharsets;
import java.security.MessageDigest;
import java.security.NoSuchAlgorithmException;
import java.util.Arrays;
import java.util.Random;

public class MainActivity extends AppCompatActivity {

    private static final String TAG = "TAG";

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);

        Button button = (Button)findViewById(R.id.submitFlag);

        button.setOnClickListener(v -> {
            EditText et = (EditText)findViewById(R.id.flag);
            String flag = et.getText().toString();
            if (flag.length() != 28) {
                Toast.makeText(MainActivity.this, "Wrong flag length", Toast.LENGTH_SHORT).show();
                return;
            }
            String result = "Wrong!";
            try {
                if (firstCheck(flag.substring(0, 9))
                        && secondCheck(flag.substring(9, 13))
                        && thirdCheck(flag.substring(13, 16))
                        && fourthCheck(flag.substring(16, 23))
                        && fifthCheck(flag.substring(23, 28))) {
                    result = "Congrats! You guessed the flag!";
                }
            } catch (NoSuchAlgorithmException e) {
                e.printStackTrace();
            }
            Toast.makeText(MainActivity.this, result, Toast.LENGTH_SHORT).show();
        });
    }

    private boolean firstCheck(String flag) {
        byte[] love = "Hackerdom".getBytes();
        byte[] check = new byte[]{4, 4, 23, 4, 38, 38, 34, 20, 30};
        byte[] inputBytes = flag.getBytes();
        for (int i = 0; i < love.length; i++) {
            love[i] = (byte)(love[i] ^ inputBytes[i]);
        }

        return Arrays.equals(check, love);
    }

    private boolean secondCheck(String flag) {
        char[] check = "0n5b".toCharArray();
        char[] flag_arr = flag.toCharArray();
        for (int i = 0; i < check.length; i++) {
            flag_arr[i] += i;
        }

        return Arrays.equals(check, flag_arr);
    }

    private boolean thirdCheck(String flag) throws NoSuchAlgorithmException {
        MessageDigest md = MessageDigest.getInstance("MD5");
        md.update(flag.getBytes(StandardCharsets.UTF_8));
        byte[] flag_hash = md.digest();
        byte[] check = new byte[]{(byte) 178, 45, (byte) 190, 74, (byte) 136, 86, (byte) 244, (byte) 152, (byte) 236, 125, 14, 102, (byte) 191, (byte) 250, 105, (byte) 203};

        return Arrays.equals(flag_hash, check);
    }

    private boolean fourthCheck(String flag) {
        Random rnd = new Random(865790124);
        long smth = rnd.nextLong();
        byte[] flag_bytes = flag.getBytes();

        long check = -1655835832096201751L;
        long val = 0;
        for (byte b: flag_bytes) {
            val = (val << 8) + (b & 0xFF);
        }
        Log.v(TAG, String.valueOf(val ^ smth));
        Log.v(TAG, String.valueOf(check));
        return (val ^ smth) == check;
    }

    private boolean fifthCheck(String flag) {
        char[] check = new char[]{51, 118, 125, 95, 114};

        for (int i = 0; i < check.length; i++) {
            if (flag.charAt(i) != check[(i + 3)  % 5])
                return false;
        }

        return true;
    }
}