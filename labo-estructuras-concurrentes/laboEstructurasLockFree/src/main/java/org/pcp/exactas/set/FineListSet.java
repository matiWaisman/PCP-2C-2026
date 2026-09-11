package org.pcp.exactas.set;

import java.util.concurrent.locks.ReentrantLock;

public class FineListSet implements ConcurrentIntSet {
    private final Node head = new Node(Long.MIN_VALUE);
    private final Node tail = new Node(Long.MAX_VALUE);

    public FineListSet() { head.next = tail; }

    @Override public boolean add(int value) {
        Node pred = head;
        Node curr = head.next;
        try{
            pred.lock.lock();
            curr.lock.lock();
            while (curr.key < value) {
                pred.lock.unlock();
                pred = curr;
                curr = curr.next;
                curr.lock.lock();
            }
            if (curr.key == value) {
                return false;
            }
            pred.next = new Node(value, curr);
            return true;
        }
        finally{
            pred.lock.unlock();
            curr.lock.unlock();
        }
        
    }

    @Override public boolean remove(int value) {
        Node pred = head;
        Node curr = head.next;
        try{
            pred.lock.lock();
            curr.lock.lock();
            while (curr.key < value) {
                pred.lock.unlock();
                pred = curr;
                curr = curr.next;
                curr.lock.lock();
            }
            if (curr.key != value) {
                return false;
            }
            pred.next = curr.next;
            return true;
        }
        finally{
            pred.lock.unlock();
            curr.lock.unlock();
        }
    }

    @Override public boolean contains(int value) {
        Node curr = head.next;
        while (curr.key < value) {
            curr = curr.next;
        }
        return curr.key == value;
    }

    @Override public boolean isImplemented() {
        return true;
    }

    private static final class Node {
        final long key;
        Node next;
        final ReentrantLock lock = new ReentrantLock();
        Node(long key) { this(key, null); } Node(long key, Node next) { this.key=key; this.next=next; }
    }
}
